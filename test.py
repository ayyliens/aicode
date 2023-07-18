# import pprint

def read_file(path):
  with open(path) as file:
    return file.read()

def test_completion():
  import openai

  msgs = [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Who won the world series in 2020?"},
    {"role": "assistant", "content": "The Los Angeles Dodgers won the World Series in 2020."},
    {"role": "user", "content": "Where was it played?"},
  ]

  compl = openai.ChatCompletion.create(
    model="gpt-3.5-turbo",
    # Not available to general public yet.
    # model="gpt-4-32k",
    # model="gpt-4",
    messages=msgs,
  )

  print(compl)

def estimate_tokens():
  estimate_tokens_from_file("...")

def estimate_tokens_from_file(src: str):
  import tiktoken
  encoding = tiktoken.get_encoding("cl100k_base")
  tokens = encoding.encode(read_file(src))
  print("len(tokens):", len(tokens))

def main():
  # print('running')
  # estimate_tokens()
  # test_completion()

if __name__ == "__main__":
  main()
