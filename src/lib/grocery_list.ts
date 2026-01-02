import z from "zod"

export const GroceryItem = z.object({
  id: z.number(),
  name: z.string(),
  quantity: z.string(),
  category: z.optional(z.string()),
  checked: z.optional(z.boolean())
})
export type GroceryItem = z.infer<typeof GroceryItem>;

export const GroceryList = z.object({
  id: z.number(),
  name: z.optional(z.string())
})
export type GroceryList = z.infer<typeof GroceryList>;
