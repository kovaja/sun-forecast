<script lang="ts">
  import type { AppLink } from '../types';
  import { Link } from 'svelte-routing';

  export let links: AppLink[];
</script>

<nav>
    {#each links as appLink}
        <Link to="{appLink.route}">
            {appLink.name}
        </Link>
    {/each}
</nav>

<style>
    nav {
        display: flex;
        align-items: center;
        justify-content: space-evenly;
        background-color: #395B64;
        /*height: 5vh;*/
    }

    /*Rest of the styles in app.css*/
</style>

