import z from "zod";

export const QuoteSchema = z.object({
  _id: z.string(),
  q: z.string(),
  a: z.string(),
  c: z.string(),
  h: z.string(),
});
