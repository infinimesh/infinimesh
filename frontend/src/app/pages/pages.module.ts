import {NgModule} from '@angular/core';

import {PagesComponent} from './pages.component';
import {DashboardModule} from './dashboard/dashboard.module';
import {PagesRoutingModule} from './pages-routing.module';
import {ThemeModule} from '../theme/theme.module';
import {MiscellaneousModule} from './miscellaneous/miscellaneous.module';
import {DeviceRegistryComponent} from './device-registry/device-registry.component';
import { DeviceDetailComponent } from './device-detail/device-detail.component';

@NgModule({
  imports: [
    PagesRoutingModule,
    ThemeModule,
    DashboardModule,
    MiscellaneousModule,
  ],
  declarations: [
    PagesComponent,
    DeviceRegistryComponent,
    DeviceDetailComponent
  ],
})
export class PagesModule {
}
