import {Injectable} from "@angular/core";
import {environment} from '../../../environments/environment';

const apiServerUrl = environment.apiServerUrl;

@Injectable()
export class UrlProviderService {

  getApiServerUrl(): string {
    if (apiServerUrl.startsWith("$")) {
      return "http://localhost:8081";
    } else {
      return apiServerUrl;
    }
  }

}
