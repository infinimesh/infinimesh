import {Injectable} from '@angular/core';
import {HttpHeaders} from '@angular/common/http';
import {NbAuthJWTToken, NbAuthService} from '@nebular/auth';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/filter';
import 'rxjs/add/operator/mergeMap';
import {UrlProviderService} from './url-provider.service';
import {HttpParamsOptions} from '@angular/common/http/src/params';

@Injectable()
export class ApiUtilService {

  private apiUrl: string;
  private httpOptions;

  constructor(private authService: NbAuthService,
              private urlProviderService: UrlProviderService) {
    this.apiUrl = urlProviderService.getApiServerUrl();
    this.authService.onTokenChange()
      .subscribe((token: NbAuthJWTToken) => {

        if (token.isValid()) {
          this.httpOptions = {
            headers: new HttpHeaders({
              'Content-Type': 'application/json',
              'Authorization': 'Bearer ' + token.getValue(),
            }),
          };
        }
      });
  }

  public getApiUrl(): string {
    return this.apiUrl;
  }

  public getHttpOptions() {
    return this.httpOptions;
  }

}
