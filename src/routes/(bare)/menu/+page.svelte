<script lang="ts">
  import MenuRenderer from "$lib/components/menuRenderer.svelte";
  // import MenuItem from "$lib/components/menuItem.svelte";

  class MenuItem {
    name: string;
    description: string;

    constructor(name: string, description: string) {
      this.name = name;
      this.description = description;
    }
  }

  const menuItems = new Map();

  menuItems.set("Mains", [
    new MenuItem("Spaghetti cacio e pepe", "Pasta prepared with a light sauce of cheese and pepper"),
    new MenuItem("Pork Braised with Chillies", "Slow-braised pork shoulder flavored with flavorful ancho chillies (not spicy)"),
    new MenuItem("Priya's Dal", "Creamy spiced Indian lentil soup"),
    new MenuItem("Pasta Bolognese", "Classic Italian pasta with red meat sauce"),
    new MenuItem("Marry me Chicken", "Oven-roasted chicken served with sundried tomatoes and a creamy herb sauce")
  ]);
  menuItems.set("Appetizers", [new MenuItem("Miso Soup", "Light miso soup served with tofu and seaweed pieces")]);
  menuItems.set("Desserts", [
    new MenuItem("Burnt Basque Cheesecake", "Fluffy cheesecake served in the Basque way"),
    new MenuItem("Ricotta Kumquat Marmalade Cake", "Rich fluffy ricotta cake served with candied Kumquats from the garden"),
    new MenuItem("Earl Grey and Apricot Hamantaschen", "Flaky bite-sized pastries baked with an earl grey and apricot filling"),
    new MenuItem("Ice Cream and Toppings", "Various flavours of ice cream served with toppings. Toppings vary with the season and contents of the fridge"),
    new MenuItem("Taiwanese Castella", "Light and airy cake served with whipped cream and berries"),
    new MenuItem("Orange and Olive Oil Upside Down Cake", "Moist and flavorful orange cake with hints of olive oil"),
    new MenuItem("Cinnamon Sugar Palmiers", "Flaky cinnamon and sugar cookies"),
    new MenuItem("Molasses Spice Cookies", "Chewy spiced molasses cookies"),
    new MenuItem("Corn Muffins", "Slightly sweet corn mffins"),
    new MenuItem("Strawberry Almond Bostock", "Fresh strawberries and almond paste baked on a slice of brioche"),
    new MenuItem("Pavlova", "Australian meringue served with seasonal fruits and whipped cream"),
  ])

</script>

<MenuRenderer menuName="Victoria Night" items={menuItems} />
