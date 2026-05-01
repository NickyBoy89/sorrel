import { z } from "zod";
import { RecipeId, Recipe, CreateRecipeInput } from "$lib/recipe";
import { Api } from "./index";

export class Recipes extends Api {
  async find(recipeId: RecipeId): Promise<Recipe> {
    return fetch(`${this.url}/api/v1/recipes/${recipeId}`).then(resp => resp.json()).then(resp => Recipe.decode(resp))
  }

  async findAllIds(): Promise<RecipeId[]> {
    return fetch(`${this.url}/api/v1/recipes`).then(resp => resp.json()).then(resp => z.array(RecipeId).decode(resp))
  }

  async create(createRecipeInput: CreateRecipeInput): Promise<RecipeId> {
    return fetch(`${this.url}/api/v1/recipes`, { method: "POST", body: JSON.stringify(createRecipeInput) }).then(resp => resp.json()).then(resp => RecipeId.decode(resp))
  }
}
