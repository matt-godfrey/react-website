import { randomQuoteQuery } from "@/api/quoteQuery";
import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { Card } from "./ui/card";

interface QuoteProps {}

export default function Quote({}: QuoteProps) {
  const { data } = useQuery(randomQuoteQuery());
  const quote = data?.q.trim();
  const author = data?.a;
  return (
    <Card>
      <div>
        <p className="font-mono p-4">&ldquo;{quote}&rdquo;</p>
        <p className="text-center">— {author}</p>
      </div>
    </Card>
  );
}
