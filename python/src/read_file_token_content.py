import zipfile
from io import StringIO
from xml.etree import ElementTree

import pdfplumber

from src.utils.utils import is_blank

# TODO: type hinting

def get_file_content(file_path, max_token_count):
  def is_docx_file(path):
    with zipfile.ZipFile(path, "r") as docx_zip:
      return "word/document.xml" in docx_zip.namelist()

  def is_text_file(path):
    try:
      with open(path, "rb") as file:
        for chunk in iter(lambda: file.read(1024), b""):
          chunk.decode("utf-8")
      return True
    except UnicodeDecodeError:
      return False

  def is_pdf_file(path):
    with open(path, "rb") as file:
      header = file.read(5)
      return header == b"%PDF-"

  if is_text_file(file_path):
    return get_processed_file_content(file_path, max_token_count, process_txt_content)
  elif is_pdf_file(file_path):
    return get_processed_file_content(file_path, max_token_count, process_pdf_content)
  elif is_docx_file(file_path):
    return get_processed_file_content(file_path, max_token_count, process_docx_content)

  return ""

def get_processed_file_content(file_path, max_token_count, process_file):
  token_count = 0
  string_builder = StringIO()

  process_file(file_path, token_count, max_token_count, string_builder)

  return string_builder.getvalue()

def process_txt_content(file_path, token_count, max_token_count, string_builder):
  # Each line ends with the "\n" character, except the last line, if the file does not end with a newline.
  with open(file_path, "r") as lines:
    for line in lines:
      token_count, is_max_token_count = process_tokens(line, string_builder, token_count, max_token_count)
      if is_max_token_count:
        return

def process_pdf_content(file_path, token_count, max_token_count, string_builder):
  with pdfplumber.open(file_path) as pdf:
    for page in pdf.pages:

      # Slicing lines costs O(n), which is why not to do that.
      lines = page.extract_text_lines()
      length_minus_one = len(lines) - 1
      for i, line in enumerate(lines):
        text = line.get("text", "")

        # This check helps not to add a "\n" character for lines without text.
        if is_blank(text):
          continue

        token_count, is_max_token_count = process_tokens(text, string_builder, token_count, max_token_count)
        if is_max_token_count:
          return
        if i < length_minus_one:
          token_count, is_max_token_count = process_newline_token(string_builder, token_count, max_token_count)
          if is_max_token_count:
            return

# TODO: WIP and check if correct, also always an extra "\n" at the end of the text
def process_docx_content(file_path, token_count, max_token_count, string_builder):
  with zipfile.ZipFile(file_path, "r") as docx_zip:
    with docx_zip.open("word/document.xml") as document_xml:
      root = ElementTree.parse(document_xml).getroot()
      namespace = {"w": "http://schemas.openxmlformats.org/wordprocessingml/2006/main"}

      for element in root.find(".//w:body", namespace):
        if element.tag == f"{{{namespace["w"]}}}p":
          texts = element.findall(".//w:t", namespace)
          for t in texts:
            token_count, is_max_token_count = process_tokens(t.text, string_builder, token_count, max_token_count)
            if is_max_token_count:
              return
          token_count, is_max_token_count = process_newline_token(string_builder, token_count, max_token_count)
          if is_max_token_count:
            return
        elif element.tag == f"{{{namespace["w"]}}}tbl":
          for row in element.findall(".//w:tr", namespace):
            for cell in row.findall(".//w:tc", namespace):
              for t in cell.findall(".//w:t", namespace):
                token_count, is_max_token_count = process_tokens(t.text, string_builder, token_count, max_token_count)
                if is_max_token_count:
                  return
                token_count, is_max_token_count = process_token(
                  "\t", string_builder, token_count, max_token_count
                )
                if is_max_token_count:
                  return
            token_count, is_max_token_count = process_newline_token(string_builder, token_count, max_token_count)
            if is_max_token_count:
              return

def process_newline_token(string_builder, token_count, max_token_count):
  return process_token("\n", string_builder, token_count, max_token_count)

def process_token(token, string_builder, token_count, max_token_count):
  string_builder.write(token)
  token_count += 1

  if token_count == max_token_count:
    return token_count, True

  return token_count, False

def process_tokens(text, string_builder, token_count, max_token_count):
  for token in basic_western_token_generator(text):
    token_count, is_max_token_count = process_token(token, string_builder, token_count, max_token_count)
    if is_max_token_count:
      return token_count, True

  return token_count, False

# This function generates tokens from Western text.
# It outputs tokens for words, whitespace characters, and punctuation marks.
# Note: This function does not support sub-word tokenization.
# For example, "unhappiness" is treated as a single token, not two tokens ("un" and "happiness").
def basic_western_token_generator(text):
  index = 0

  def is_space_or_punctuation(c):
    return c.isspace() or c in {",", ".", "?", "!", ";", ":", "(", ")", "[", "]"}

  while index < len(text):
    if is_space_or_punctuation(text[index]):
      yield text[index]
      index += 1
    else:
      string_builder = StringIO()
      while index < len(text) and not is_space_or_punctuation(text[index]):
        string_builder.write(text[index])
        index += 1
      yield string_builder.getvalue()
