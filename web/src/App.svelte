<script lang="ts">
  import { Router, Link, Route } from "svelte-routing";
  import Events from './routes/Events.svelte';
  import Forecast from './routes/Forecast.svelte';
  import type { AppLink } from './lib/types';
  import Tabs from './lib/components/Tabs.svelte';
  import Diffs from './routes/Diffs.svelte';
  import Modal from './lib/components/Modal/Modal.svelte';
  import Dashboard from './routes/Dashboard.svelte';

  const links: AppLink[] = [
    {
      name: 'Forecast',
      route: '/'
    },
    {
      name: 'Dashboard',
      route: '/dashboard'
    },
    {
      name: 'Diffs',
      route: '/diffs'
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
      <Route path="/diffs" component={Diffs}/>
      <Route path="/dashboard" component={Dashboard}/>
      <Route path="/" component={Forecast}/>
    </div>
    <Modal />
  </Router>
</main>

<style>
  main {
    width: 100%;
    height: 100%;
  }

  .route-main {
    /* 100% - tabs height */
    height: calc(100% - 30px);
  }
</style>
