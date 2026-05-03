import z from "zod";
import { GroceryItem } from "./grocery_list";

export const RecipeId = z.int().brand<"RecipeId">()
export type RecipeId = z.TypeOf<typeof RecipeId>;

export const Recipe = z.object({
  id: RecipeId,
  name: z.string(),
  color: z.string(),
  ingredients: z.array(GroceryItem),
});
export type Recipe = z.TypeOf<typeof Recipe>;

const { id: _, ...createRecipeInput } = Recipe.shape;
export const CreateRecipeInput = z.object({
  ...createRecipeInput
})
export type CreateRecipeInput = z.TypeOf<typeof CreateRecipeInput>;
