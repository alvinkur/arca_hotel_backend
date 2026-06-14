"""Flask microservice serving TF-IDF room recommendations."""

import os
import logging

import requests
from flask import Flask, jsonify, request

from trainer import recommend, train

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger(__name__)

app = Flask(__name__)

ROOM_SERVICE_URL = os.getenv("ROOM_SERVICE_URL", "http://localhost:8002")

_vectorizer = None
_df = None
_tfidf_matrix = None
_rooms = []


def fetch_room_types() -> list[dict]:
    try:
        resp = requests.get(f"{ROOM_SERVICE_URL}/api/room-types", timeout=5)
        resp.raise_for_status()
        data = resp.json()
        logger.info("Fetched %d room types from room-service", len(data))
        return data
    except Exception as e:
        logger.error("Failed to fetch room types: %s", e)
        return []


def fetch_rooms() -> list[dict]:
    try:
        resp = requests.get(f"{ROOM_SERVICE_URL}/api/rooms", timeout=5)
        resp.raise_for_status()
        data = resp.json()
        logger.info("Fetched %d rooms from room-service", len(data))
        return data
    except Exception as e:
        logger.error("Failed to fetch rooms: %s", e)
        return []


def init_model():
    global _vectorizer, _df, _tfidf_matrix, _rooms

    room_types = fetch_room_types()
    if not room_types:
        logger.warning("No room types available — model not initialized")
        return

    _rooms = fetch_rooms()
    _vectorizer, _df, _tfidf_matrix = train(room_types)
    logger.info("TF-IDF model trained on %d room types, %d rooms loaded",
                len(_df), len(_rooms))


@app.route("/health", methods=["GET"])
def health():
    if _vectorizer is None:
        return jsonify({"status": "not ready", "model_loaded": False}), 503
    return jsonify({
        "status": "ok",
        "model_loaded": True,
        "room_type_count": len(_df),
        "room_count": len(_rooms),
    })


@app.route("/recommend", methods=["POST"])
def recommend_endpoint():
    if _vectorizer is None:
        return jsonify({"error": "Model not initialized — no room types available"}), 503

    body = request.get_json(silent=True)
    if not body or "message" not in body:
        return jsonify({"error": "Missing 'message' field in request body"}), 400

    message = body["message"].strip()
    if not message:
        return jsonify({"error": "Message must not be empty"}), 400

    results = recommend(message, _vectorizer, _df, _tfidf_matrix, _rooms)
    return jsonify(results)


if __name__ == "__main__":
    port = int(os.getenv("ML_SERVICE_PORT", "8010"))
    init_model()
    logger.info("ML Service running on :%d", port)
    app.run(host="0.0.0.0", port=port, debug=False)
