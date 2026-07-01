import { Card } from "../ui/card";
import { Skeleton } from "../ui/skeleton";

interface QuoteSkeletonProps {}

export default function QuoteSkeleton({}: QuoteSkeletonProps) {
  return (
    <Card>
      <div className="flex flex-col gap-2 items-center w-200 h-15">
        <Skeleton className="w-100 h-10 rounded-full"></Skeleton>
        <Skeleton className="w-20 h-8 rounded-full text-center"> </Skeleton>
      </div>
    </Card>
  );
}
