<script lang="ts">
    let { text, href, color = "#ffffff", bgColor = "#212121" } = $props();
</script>

<div class="button-container rounded-xl text-center">
    <a class="shadowed-button shadowed-button-container px-3 py-1.5 rounded-sm cursor-pointer" style="--button-bg-color: {bgColor}; --button-fg-color: {color}" href={href}>{text}</a>
</div>

<style>
    .shadowed-button {
        position:relative;
        transition:all .3s ease;
    }
    .shadowed-button:focus {
       outline:none;
    }
    .shadowed-button:focus:before,
    .shadowed-button:hover:before {
        left:0;
        top:0;
    }
    .shadowed-button:before {
        content:"";
        height:100%;
        position:absolute;
        transition:all .3s ease-in-out;
        width:100%;
        z-index:-1
    }
    .shadowed-button-container {
        background:var(--button-bg-color);
        border:1px solid var(--button-fg-color);
        border-radius:9px;
        padding:.25rem .5rem
    }
    .shadowed-button-container:focus,
    .shadowed-button-container:hover {
        color:inherit;
        outline:none
    }
    .shadowed-button-container:before {
        border:1px solid var(--button-fg-color);
        border-radius:9px;
        left:4px;
        top:4px
    }
</style>