
const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

// Retrieve all news articles from Go backend.
export async function getNews() {
    const response = await fetch(`${API_URL}/news`);

    if (!response.ok) {
        throw new Error("Failed to fetch news");
    }

    return await response.json();
}
