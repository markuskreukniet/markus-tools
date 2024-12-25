import zipfile
from io import StringIO

import pdfplumber

from src.utils.utils import is_blank

# TODO: use python-docx

# TODO: duplicate code in get_pdf_content and get_txt_content

# TODO: check functions: get_file_content, is_text_file, and is_pdf_file
def get_file_content(file_path, max_token_count):
  def is_docx_file(path):
    try:
      with zipfile.ZipFile(path, "r") as docx_zip:
        return "word/document.xml" in docx_zip.namelist()
    except zipfile.BadZipFile:
      return False

  def is_text_file(path, chunk_size=1024):
    try:
      with open(path, "rb") as file:
        for chunk in iter(lambda: file.read(1024), b""):
          chunk.decode("utf-8")
      return True
    except (UnicodeDecodeError, OSError):
      return False

  def is_pdf_file(path):
    try:
      with open(path, "rb") as file:
        header = file.read(5)
        return header == b"%PDF-"
    except Exception:
      return False

  if is_text_file(file_path):
    return get_txt_content(file_path, max_token_count)
  elif is_pdf_file(file_path):
    return get_pdf_content(file_path, max_token_count)

def get_pdf_content(file_path, max_token_count):
  token_count = 0
  string_builder = StringIO()

  with pdfplumber.open(file_path) as pdf:
    for page in pdf.pages:

      # Slicing lines costs O(n), which is why not to do that.
      lines = page.extract_text_lines()
      length_minus_one = len(lines) - 1
      for i, line in enumerate(lines):
        text = line.get("text", "").strip()
        if is_blank(text):
          continue

        for token in basic_western_token_generator(text):
          string_builder.write(token)
          token_count += 1
          if token_count == max_token_count:
            return string_builder.getvalue()
        if i < length_minus_one:
          string_builder.write("\n")
          token_count += 1
          if token_count == max_token_count:
            return string_builder.getvalue()

  return string_builder.getvalue()

def get_txt_content(file_path, max_token_count):
  token_count = 0
  string_builder = StringIO()

  # Each line ends with the "\n" character, except the last line, if the file does not end with a newline.
  with open(file_path, "r") as lines:
    for line in lines:
      for token in basic_western_token_generator(line):
        string_builder.write(token)
        token_count += 1
        if token_count == max_token_count:
          return string_builder.getvalue()

  return string_builder.getvalue()

# This function generates tokens from Western text.
# It outputs tokens for words, whitespace characters, and punctuation marks.
# Note: This function does not support sub-word tokenization.
# For example, "unhappiness" is treated as a single token, not two tokens ("un" and "happiness").
def basic_western_token_generator(text):
  index = 0

  def is_space_or_punctuation(c):
    return c.isspace() or c in {',', '.', '?', '!', ';', ':', '(', ')', '[', ']'}

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
