import requests

GOOGLE_API_KEY = "AIzaSyDIx-yxvJ8mjVETFLikuWO6XCjYnmUR-q4"


def callLLM(prompt):
    response = requests.post(
        url="https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-lite:generateContent",
        headers={
            "x-goog-api-key": GOOGLE_API_KEY,
            "Content-Type": "application/json",
        },
        json={"contents": [{"parts": [{"text": prompt}]}]},
    )
    return response.json()


if __name__ == "__main__":
    prompt = "Explain the theory of relativity in a few words"
    response = callLLM(prompt)
    print(response)

