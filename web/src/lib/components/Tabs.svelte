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


    /*Rest of the styles in app.css*/
</style>

