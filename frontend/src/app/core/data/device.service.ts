import {Observable} from 'rxjs';
import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {NbAuthJWTToken, NbAuthService} from '@nebular/auth';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/filter';
import 'rxjs/add/operator/mergeMap';
import {UrlProviderService} from './url-provider.service';
import {ApiUtilService} from './api-util.service';
import {NamespaceService} from './namespace.service';

@Injectable()
export class DeviceService {

  constructor(private http: HttpClient,
              private apiUtilService: ApiUtilService,
              private namespaceService: NamespaceService) {

  }

  getAll(): Observable<any> {
    const namespace = this.namespaceService.getSelected();
    const url = `${this.apiUtilService.getApiUrl()}/devices?namespace=${namespace.name}`;
    return this.http.get(url, this.apiUtilService.getHttpOptions())
      .map((response: any) => response.devices);
  }

}
