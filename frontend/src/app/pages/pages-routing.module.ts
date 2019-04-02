import {RouterModule, Routes} from '@angular/router';
import {NgModule} from '@angular/core';

import {PagesComponent} from './pages.component';
import {DashboardComponent} from './dashboard/dashboard.component';
import {DeviceRegistryComponent} from "./device-registry/device-registry.component";
import {DeviceDetailComponent} from "./device-detail/device-detail.component";

const routes: Routes = [{
  path: '',
  component: PagesComponent,
  children: [
    {
      path: 'dashboard',
      component: DashboardComponent,
    },
    {
      path: 'devices',
      component: DeviceRegistryComponent,
    },
    {
      path: 'devices/:deviceId',
      component: DeviceDetailComponent,
    },
    {
      path: '',
      redirectTo: 'devices', /* default view */
      pathMatch: 'full',
    },
  ],
}];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class PagesRoutingModule {
}
