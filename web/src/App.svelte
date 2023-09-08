<script lang="ts">
  import { Router, Link, Route } from "svelte-routing";
  import Events from './routes/Events.svelte';
  import Forecast from './routes/Forecast.svelte';
  import type { AppLink } from './lib/types';
  import Tabs from './lib/components/Tabs.svelte';

  const links: AppLink[] = [
    {
      name: 'Forecast',
      route: '/'
    },
    {
      name: 'Events',
      route: '/events'
    }
  ]

  let url = window.location.pathname;
</script>

<main>
    <Router {url}>
        <Tabs links={links}/>
        <div class="route-main">
            <Route path="/events" component={Events}/>
            <Route path="/">
                <Forecast/>
            </Route>
        </div>
    </Router>
</main>

<style>
    main {
        width: 100%;
        height: 100%;
    }
</style>
