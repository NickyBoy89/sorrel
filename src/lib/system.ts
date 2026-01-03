import z from "zod";

export const User = z.object({
  id: z.number(),
  display_name: z.string()
})
export type User = z.infer<typeof User>;
