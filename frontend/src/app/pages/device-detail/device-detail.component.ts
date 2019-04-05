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
  stateUpdated = false;
  JSON = JSON;

  constructor(private route: ActivatedRoute,
              private deviceService: DeviceService) { }

  ngOnInit() {
    this.route.paramMap.subscribe(params => {
      const deviceId = params.get("deviceId")
      this.device$ = this.deviceService.getOne(deviceId);
      this.deviceService.getState(deviceId).subscribe(state => {
        if(state) {
          this.state = state;
        } else {
          this.state = {
            reported: {}
          };
        }
        this.stateSubscription = this.deviceService.streamState(deviceId).subscribe((data) => {
          if (data.reportedState) {
            this.state.reported.data = data.reportedState.data;
            this.state.reported.timestamp = data.reportedState.timestamp;
            this.state.reported.version = data.reportedState.version;
            this.stateUpdated = true
          }
          setTimeout(() => {
            this.stateUpdated = false
          }, 500);
        });
      })
    })
  }

  ngOnDestroy(): void {
    this.stateSubscription.unsubscribe();
  }

}
