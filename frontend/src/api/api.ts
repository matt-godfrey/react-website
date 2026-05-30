// https://api.mattgodfrey.xyz/api/random

import { QuoteSchema } from "../schemas/quoteSchema";

export async function getRandomQuote() {
  const response = await fetch("https://api.mattgodfrey.xyz/api/random");
  const data = await response.json();
  return QuoteSchema.parse(data);
}
