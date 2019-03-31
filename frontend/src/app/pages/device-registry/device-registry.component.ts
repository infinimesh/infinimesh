import {Component, OnInit} from '@angular/core';
import {DeviceService} from '../../core/data/device.service';
import {Observable} from 'rxjs';
import {NamespaceService} from '../../core/data/namespace.service';

@Component({
  selector: 'device-registry',
  templateUrl: './device-registry.component.html',
  styleUrls: ['./device-registry.component.scss']
})
export class DeviceRegistryComponent implements OnInit {

  public devices$: Observable<any>;

  constructor(private deviceService: DeviceService,
              private namespaceService: NamespaceService) {
    this.namespaceService.selectedChange.subscribe(() => this.requestDevices());
  }

  ngOnInit() {
    this.requestDevices();
  }

  requestDevices() {
    if (this.namespaceService.getSelected() !== undefined) {
      this.devices$ = this.deviceService.getAll();
    }
  }
}
