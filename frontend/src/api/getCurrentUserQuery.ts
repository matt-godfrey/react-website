import { getCurrentUser } from "./api";

export function getCurrentUserQuery() {
  return {
    queryKey: ["me"],
    queryFn: getCurrentUser,
    retry: false,
  };
}
