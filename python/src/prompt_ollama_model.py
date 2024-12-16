import http.client
import json

file_name = "\"fileName\""

# TODO: error handling
def prompt_ollama_model():
    prompt = f"The content of a text file follows. Give a good file name for that file in a JSON format. The file name should be the property's {file_name} value as one string.\n\nThe cat (Felis catus), or domestic cat, is a small carnivorous mammal and the only domesticated species in the Felidae family. Domesticated around 7500 BC in the Near East, cats are valued as pets and for controlling vermin. They are agile hunters with retractable claws, sharp teeth, excellent night vision, and a keen sense of smell. Though social, cats hunt alone, often at dawn and dusk. They communicate through vocalizations (meowing, purring, hissing) and body language, can hear high-frequency sounds, and use pheromones for signaling."

    connection = http.client.HTTPConnection("localhost:11434")

    connection.request(
        method="POST",
        url="/api/generate",
        headers={
            "Content-Type": "application/json"
        },

        body=json.dumps({
            "model": "llama3.2:3B",
            "prompt": prompt,
            "stream": False,
            "temperature": 0.1,
            "top_p": 0.6,

            # defaults:
            "mirostat": 0,
            "mirostat_eta": 0.1,
            "mirostat_tau": 0.5,
            "repeat_last_n": 64,
            "repeat_penalty": 1.1,
            "num_predict": -1,
            "tfs_z": 1,
            "seed": 0,
            "top_k": 40,
            "num_ctx": 2048,
            "min_p": 0.0
        }),
    )

    result = ""
    response = connection.getresponse()

    if response.status == 200:
        result = json.loads(response.read().decode('utf-8'))["response"]
    else:
        print("Error")

    connection.close()

    return result

def create_file_name():
  prompt_response = prompt_ollama_model()
  index = prompt_response.find(file_name)

  if index == -1:
    return ""

  return prompt_response[index:]

print(create_file_name())
