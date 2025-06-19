import json
import sys
import numpy as np
import torch
from sentence_transformers import SentenceTransformer, util
import os
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'

device = torch.device('cuda') if torch.cuda.is_available() else torch.device('cpu')

input_data = json.load(sys.stdin)
field_descriptions = input_data["fieldDescriptions"]
api_methods = input_data["apiMethods"]

model = SentenceTransformer('sentence-transformers/stsb-roberta-large').to(device)

fields = list(field_descriptions.keys())
descriptions = list(field_descriptions.values())
api_descriptions = [api["Description"] for api in api_methods]

description_embeddings = model.encode(descriptions, convert_to_tensor=True, device=device).cpu().numpy()
api_embeddings = model.encode(api_descriptions, convert_to_tensor=True, device=device).cpu().numpy()

cosine_similarities = util.pytorch_cos_sim(description_embeddings, api_embeddings).cpu().numpy()

threshold = 0.75 
recommended_apis = {}
for i, field in enumerate(fields):
    best_match_idx = np.argmax(cosine_similarities[i])
    best_match_score = float(cosine_similarities[i][best_match_idx])
    if best_match_score >= threshold:
        api_name = api_methods[best_match_idx]["Name"]
        if api_name not in recommended_apis:
            recommended_apis[api_name] = []
        recommended_apis[api_name].append({"field": field, "score": best_match_score})

for api in recommended_apis:
    unique_fields = {}
    for item in recommended_apis[api]:
        field = item["field"]
        score = item["score"]
        if field not in unique_fields:
            unique_fields[field] = score
        elif score > unique_fields[field]:
            unique_fields[field] = score
    recommended_apis[api] = [{"field": field, "score": score} for field, score in unique_fields.items()]

print(json.dumps(recommended_apis, indent=2))