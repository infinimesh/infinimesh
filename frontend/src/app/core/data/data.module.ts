import {ModuleWithProviders, NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';

import {UserService} from './user.service';
import {UrlProviderService} from "./url-provider.service";
import {DeviceService} from "./device.service";
import {ApiUtilService} from "./api-util.service";
import {NamespaceService} from "./namespace.service";

const SERVICES = [
  UserService,
  DeviceService,
  UrlProviderService,
  ApiUtilService,
  NamespaceService
];

@NgModule({
  imports: [
    CommonModule,
  ],
  providers: [
    ...SERVICES,
  ],
})
export class DataModule {
  static forRoot(): ModuleWithProviders {
    return <ModuleWithProviders>{
      ngModule: DataModule,
      providers: [
        ...SERVICES,
      ]
    };
  }
}
