import { randomQuoteQuery } from "@/api/quoteQuery";
import { useQuery } from "@tanstack/react-query";
import { Card } from "./ui/card";
import QuoteSkeleton from "./skeletons/quoteSkeleton";

interface QuoteProps {}

export default function Quote({}: QuoteProps) {
  const { data, isLoading } = useQuery(randomQuoteQuery());
  const quote = data?.q.trim();
  const author = data?.a;
  return isLoading ? (
    // <p>Loading...</p>
    <QuoteSkeleton />
  ) : (
    <Card>
      <div>
        <p className="font-mono p-4">&ldquo;{quote}&rdquo;</p>
        <p className="text-center">— {author}</p>
      </div>
    </Card>
  );
}
