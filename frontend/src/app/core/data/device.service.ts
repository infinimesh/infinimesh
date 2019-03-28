import {Observable} from 'rxjs';
import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {NbAuthJWTToken, NbAuthService} from "@nebular/auth";
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/filter';
import 'rxjs/add/operator/mergeMap';
import {UrlProviderService} from "./url-provider.service";

@Injectable()
export class DeviceService {

  private apiUrl: string;
  private httpOptions;

  constructor(private http: HttpClient,
              private authService: NbAuthService,
              private urlProviderService: UrlProviderService) {
    this.apiUrl = urlProviderService.getApiServerUrl();
    this.authService.onTokenChange()
      .subscribe((token: NbAuthJWTToken) => {

        if (token.isValid()) {
          this.httpOptions = {
            headers: new HttpHeaders({
              'Content-Type': 'application/json',
              'Authorization': 'Bearer ' + token.getValue()
            })
          };
        }

      });
  }

  getAll(): Observable<any> {
    return this.http.get(`${this.apiUrl}/devices?namespace=shared-project`, this.httpOptions)
      .map((response: any) => response.devices);
  }

}
