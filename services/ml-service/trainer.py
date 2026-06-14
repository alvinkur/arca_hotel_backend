"""TF-IDF + Cosine Similarity room recommendation engine."""

import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity

INDONESIAN_STOP_WORDS = [
    "saya", "ingin", "mau", "cari", "butuh", "yang", "dan", "atau", "di", "ke",
    "buat", "untuk", "dengan", "bisa", "ada", "kamar", "hotel", "tidak", "juga",
    "tolong", "si", "nih", "dong", "ya", "gak", "nggak", "aja", "deh", "kok",
    "itu", "ini", "aku", "gua", "gue", "the", "a", "an", "is", "are", "for",
    "kami", "rekomendasi", "sarankan", "pilih", "pilihan", "cocok", "bagus",
]


def train(room_types: list[dict]) -> tuple[TfidfVectorizer, pd.DataFrame]:
    """Train TF-IDF vectorizer on room type name + description.

    Args:
        room_types: list of dicts with keys: name, price, description

    Returns:
        (fitted vectorizer, TF-IDF matrix as DataFrame)
    """
    df = pd.DataFrame(room_types)
    df["features"] = df["name"] + " " + df["description"]

    vectorizer = TfidfVectorizer(stop_words=INDONESIAN_STOP_WORDS)
    tfidf_matrix = vectorizer.fit_transform(df["features"])

    return vectorizer, df, tfidf_matrix


def recommend(
    query: str,
    vectorizer: TfidfVectorizer,
    df: pd.DataFrame,
    tfidf_matrix,
    rooms: list[dict] = None,
) -> list[dict]:
    """Rank room types by cosine similarity to the query.

    Args:
        query: user's natural language preference message
        vectorizer: fitted TfidfVectorizer
        df: DataFrame with room type data (name, price, description)
        tfidf_matrix: fitted TF-IDF matrix
        rooms: list of room dicts with room_number, id_room_type, availability

    Returns:
        list of dicts sorted by score descending, each with:
        name, price, description, score, room_numbers
    """
    query_vec = vectorizer.transform([query])
    similarities = cosine_similarity(query_vec, tfidf_matrix).flatten()

    room_map = {}
    if rooms:
        for r in rooms:
            if r.get("availability", True):
                rt_id = r.get("id_room_type")
                if rt_id not in room_map:
                    room_map[rt_id] = []
                room_map[rt_id].append(r["room_number"])

    results = []
    for i, score in enumerate(similarities):
        rt_id = int(df.iloc[i]["id_room_type"])
        results.append({
            "name": df.iloc[i]["name"],
            "price": float(df.iloc[i]["price"]),
            "description": df.iloc[i]["description"],
            "score": round(float(score), 4),
            "room_numbers": sorted(room_map.get(rt_id, [])),
        })

    results.sort(key=lambda r: r["score"], reverse=True)
    return results
