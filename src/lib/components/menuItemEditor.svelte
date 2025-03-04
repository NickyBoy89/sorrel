<script lang="ts"></script>

<div class="menu-item-editor rounded-md p-4">
    <form action="" class="flex flex-col">
        <div class="grid grid-cols-5">
            <select name="" id="" class="rounded-sm menu-section-selector pl-2">
                <option value="mains">Mains</option>
                <option value="desserts">Desserts</option>
            </select>
            <div class="col-start-5 text-right">Not saved</div>
        </div>
        <!-- <label for="edit-item-name">Name</label> -->
        <input type="text" id="edit-item-name" name="edit-item-name" class="item-edit-input rounded-sm mt-2 pl-2" placeholder="Name of dish...">
        <!-- <label for="edit-item-desc">Description</label> -->
        <input type="text" id="edit-item-desc" class="item-edit-input rounded-sm mt-2 pl-2" placeholder="Description (Optional)">
    </form>
</div>

<style>
    .menu-item-editor {
        background-color: #665c54;
    }

    .menu-section-selector {
        background-color: #3c3836;
    }

    .item-edit-input {
        background-color: #3c3836;
        height: 2rem;
    }
</style>