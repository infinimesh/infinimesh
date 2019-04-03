import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {DeviceService} from "../../core/data/device.service";

@Component({
  selector: 'device-detail',
  templateUrl: './device-detail.component.html',
  styleUrls: ['./device-detail.component.scss']
})
export class DeviceDetailComponent implements OnInit, OnDestroy {

  device$;
  state;
  stateSubscription;
  JSON = JSON;

  constructor(private route: ActivatedRoute,
              private deviceService: DeviceService) { }

  ngOnInit() {
    this.route.paramMap.subscribe(params => {
      const deviceId = params.get("deviceId")
      this.device$ = this.deviceService.getOne(deviceId);
      this.deviceService.getState(deviceId).subscribe(state => {
        this.state = state;
        this.stateSubscription = this.deviceService.streamState(deviceId).subscribe((data) => {
          console.log(data);
          this.state.reported.data = data.reportedDelta.data;
          this.state.reported.timestamp = data.reportedDelta.timestamp;
          this.state.reported.version = data.reportedDelta.version;
        });
      })
    })
  }

  ngOnDestroy(): void {
    this.stateSubscription.unsubscribe();
  }

}
