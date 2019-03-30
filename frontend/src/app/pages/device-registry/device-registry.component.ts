import {Component, OnInit} from '@angular/core';
import {DeviceService} from "../../core/data/device.service";
import {Observable} from "rxjs";

@Component({
  selector: 'device-registry',
  templateUrl: './device-registry.component.html',
  styleUrls: ['./device-registry.component.scss']
})
export class DeviceRegistryComponent implements OnInit {

  public devices$: Observable<any>;

  constructor(private deviceService: DeviceService) {
  }

  ngOnInit() {
    this.devices$ = this.deviceService.getAll();
  }
}
