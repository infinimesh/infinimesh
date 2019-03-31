import {ModuleWithProviders, NgModule, Optional, SkipSelf} from '@angular/core';
import {CommonModule} from '@angular/common';
import {NbAuthJWTToken, NbAuthModule, NbPasswordAuthStrategy} from '@nebular/auth';

import {throwIfAlreadyLoaded} from './module-import-guard';
import {DataModule} from './data/data.module';
import {AnalyticsService} from './utils/analytics.service';
import {UserNameLoginComponent} from "./username-login/user-name-login.component";
import {NbAlertModule, NbButtonModule, NbCheckboxModule, NbInputModule} from "@nebular/theme";
import {FormsModule} from "@angular/forms";
import {RouterModule} from "@angular/router";
import {AuthGuard} from "./auth-gard.service";
import {HttpClientModule} from "@angular/common/http";

import {environment} from '../../environments/environment';

// use localhost, if no server url is configured
let apiServerUrl = environment.apiServerUrl;
if (apiServerUrl.startsWith("$")) {
  apiServerUrl = "http://localhost:8081";
}

export const NB_CORE_PROVIDERS = [
  AuthGuard,
  ...DataModule.forRoot().providers,
  ...NbAuthModule.forRoot({

    strategies: [
      NbPasswordAuthStrategy.setup({
        name: 'email',
        token: {
          class: NbAuthJWTToken,
          key: 'token'
        },
        baseEndpoint: apiServerUrl,
        login: {
          endpoint: '/account/token',
          method: 'post'
        }
      }),
    ],
    forms: {},
  }).providers,
  AnalyticsService,
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    NbAlertModule,
    NbInputModule,
    NbButtonModule,
    NbCheckboxModule,
    NbAuthModule,
    HttpClientModule
  ],
  exports: [
    NbAuthModule,
  ],
  declarations: [
    UserNameLoginComponent
  ],
})
export class CoreModule {
  constructor(@Optional() @SkipSelf() parentModule: CoreModule) {
    throwIfAlreadyLoaded(parentModule, 'CoreModule');
  }

  static forRoot(): ModuleWithProviders {
    return <ModuleWithProviders>{
      ngModule: CoreModule,
      providers: [
        ...NB_CORE_PROVIDERS,
      ],
    };
  }
}
