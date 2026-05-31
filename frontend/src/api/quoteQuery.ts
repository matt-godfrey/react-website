import { getRandomQuote } from "./api";

export function randomQuoteQuery() {
  return {
    queryKey: ["random-quote"],
    queryFn: () => getRandomQuote(),
    staleTime: 1000 * 60 * 10,
    refetchOnWindowFocus: false,
  };
}
