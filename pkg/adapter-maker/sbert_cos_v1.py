import os
import json
import sys
import numpy as np
import torch
import torch.nn as nn
from sentence_transformers import SentenceTransformer
device = torch.device('cuda') if torch.cuda.is_available() else torch.device('cpu')
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'

class BertSimilarityModel(nn.Module):
    def __init__(self, bert_model_name='bert-base-nli-mean-tokens'):
        super(BertSimilarityModel, self).__init__()
        self.bert_model = SentenceTransformer(bert_model_name)
        self.linear1 = nn.Linear(768, 512) 
        self.linear2 = nn.Linear(512, 256)
        self.linear3 = nn.Linear(256, 128)
        self.cosine = nn.CosineSimilarity(dim=1)

    def forward(self, desc1_embedding, desc2_embedding):
        # extract feature
        desc1_features = torch.relu(self.linear1(desc1_embedding))
        desc1_features = torch.relu(self.linear2(desc1_features))
        desc1_features = torch.relu(self.linear3(desc1_features))

        desc2_features = torch.relu(self.linear1(desc2_embedding))
        desc2_features = torch.relu(self.linear2(desc2_features))
        desc2_features = torch.relu(self.linear3(desc2_features))

        # cosine similarity
        similarity = self.cosine(desc1_features, desc2_features)
        return similarity

model = BertSimilarityModel()

input_data = json.load(sys.stdin)
field_descriptions = input_data["fieldDescriptions"]
api_methods = input_data["apiMethods"]

fields = list(field_descriptions.keys())
descriptions = list(field_descriptions.values())
api_descriptions = [api["Description"] for api in api_methods]

if not fields or not api_methods:
    raise ValueError("Field descriptions or API methods are empty.")

description_embeddings = model.bert_model.encode(descriptions, convert_to_tensor=True, device=device)
api_embeddings = model.bert_model.encode(api_descriptions, convert_to_tensor=True, device=device)

if len(description_embeddings) == 0 or len(api_embeddings) == 0:
    raise ValueError("Embedding generation failed for descriptions or API methods.")

cosine_similarities = np.zeros((len(description_embeddings), len(api_embeddings)))
for i, field_desc_embedding in enumerate(description_embeddings):
    for j, api_desc_embedding in enumerate(api_embeddings):
        cosine_similarities[i][j] = model(field_desc_embedding.unsqueeze(0), api_desc_embedding.unsqueeze(0)).item()

threshold = 0.80 
top_k = 3 

recommended_apis = {}
for i, field in enumerate(fields):
    if len(cosine_similarities[i]) == 0:
        continue  
    best_match_idxs = np.argsort(-cosine_similarities[i])[:top_k]  # Top k recommendations
    for best_match_idx in best_match_idxs:
        best_match_score = float(cosine_similarities[i][best_match_idx])
        if best_match_score >= threshold:
            api_name = api_methods[best_match_idx]["Name"]
            if api_name in recommended_apis:
                recommended_apis[api_name].append({"field": field, "score": best_match_score})
            else:
                recommended_apis[api_name] = [{"field": field, "score": best_match_score}]

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
