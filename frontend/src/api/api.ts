// https://api.mattgodfrey.xyz/api/random

import { QuoteSchema } from "../schemas/quoteSchema";

export async function getRandomQuote() {
  const response = await fetch("https://api.mattgodfrey.xyz/api/random");
  const data = await response.json();
  return QuoteSchema.parse(data);
}

export async function getCurrentUser() {
  const API_URL = import.meta.env.VITE_API_URL;
  const response = await fetch(`${API_URL}/auth/me`, {
    credentials: "include",
  });
  // const data = await response.json();

  if (response.status === 401) {
    return null;
  }
  if (!response.ok) {
    throw new Error("Failed to fetch current user");
  }
  return response.json();
}
