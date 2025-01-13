from dataclasses import dataclass
from pathlib import Path
from typing import Optional

from src.prompt_llama_cpp_model import prompt_llama_cpp_model
from src.prompt_ollama_model import prompt_ollama_model
from src.read_file_token_content import basic_western_token_generator, get_file_content
from src.utils.utils import is_blank, FILE_NAME, INSTRUCTION, NUMBER_OF_CONTEXT_TOKENS

@dataclass
class PromptModel:
  content: str
  number_of_input_tokens: int
  model_file_path: Optional[str] = None

def create_file_name_without_extension(prompt_model):
  if prompt_model.model_file_path is None:
    prompt_response = prompt_ollama_model(prompt_model.content, prompt_model.number_of_input_tokens)
  else:
    prompt_response = prompt_llama_cpp_model(
      prompt_model.content, prompt_model.number_of_input_tokens, prompt_model.model_file_path
    )

  if prompt_response == "":
    return ""

  index = prompt_response.find(FILE_NAME)

  if index == -1:
    return ""

  def slice_and_find_index(response, start_index):
    response = response[start_index:]
    return response, response.find("\"")

  prompt_response, index = slice_and_find_index(prompt_response, index + len(FILE_NAME))
  prompt_response, index = slice_and_find_index(prompt_response, index+1)

  return prompt_response[:index]

def change_file_name(path, prompt_model):
  name = create_file_name_without_extension(prompt_model)

  if is_blank(name):
    return

  new_file_name = name + "".join(path.suffixes)

  path.rename(path.parent / new_file_name)

def approximate_western_token_count(text):
  token_count = 0

  for _ in basic_western_token_generator(text):
    token_count += 1

  return token_count

def change_file_name_by_content(file_path, model_file_path = None):
  path = Path(file_path)

  if not path.is_file() and path.stat().st_size > 0:
    return

  number_of_input_tokens = round(NUMBER_OF_CONTEXT_TOKENS * 0.6)
  content = get_file_content(file_path, number_of_input_tokens - approximate_western_token_count(INSTRUCTION))

  if is_blank(content):
    return

  model = PromptModel(
    content=content,
    number_of_input_tokens=number_of_input_tokens
  )

  if model_file_path is not None:
    model.model_file_path = model_file_path

  change_file_name(path, model)

change_file_name_by_content(
  "C:\\Users\\testUser\\Desktop\\test\\test.txt",
  r"C:\Users\testUser\Downloads\bartowski Llama-3.2-3B-Instruct-Q4_K_M.gguf"
)
