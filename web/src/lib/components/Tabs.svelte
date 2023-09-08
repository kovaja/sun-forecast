<script lang="ts">
  import type { AppLink } from '../types';
  import { Link } from 'svelte-routing';

  export let links: AppLink[];

  let tabsOpened = false

  function toggleTabs (e: MouseEvent) {
    e.preventDefault()
    tabsOpened = true
  }

  function closeTabs () {
    tabsOpened = false
  }
</script>

<nav class={'tabs-full' + (tabsOpened ? ' tabs-opened' : '')}>
        {#each links as appLink}
            <Link to="{appLink.route}" on:click={closeTabs}>
                {appLink.name}
            </Link>
        {/each}
</nav>

<nav class={'tabs-mobile-toggle' + (tabsOpened ? ' tabs-opened' : '')}>
    <button on:click={toggleTabs}>
        Menu
    </button>
</nav>

<style>
    .tabs-full {
        display: none;
        &.tabs-opened {
            display: flex;
        }
    }

    .tabs-mobile-toggle.tabs-opened {
        display: none;
    }

    .tabs-mobile-toggle button {
        font-size: 8px;
        padding: 0;
        border: none;
        text-align: center;
        width: 100%;
    }

    @media (min-width: 480px) {
        .tabs-mobile-toggle {
            display: none;
        }
        .tabs-full {
            display: flex;
        }
    }

    /*Rest of the styles in app.css*/
</style>

