import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {DeviceService} from '../../core/data/device.service';

@Component({
  selector: 'device-detail',
  templateUrl: './device-detail.component.html',
  styleUrls: ['./device-detail.component.scss']
})
export class DeviceDetailComponent implements OnInit, OnDestroy {

  deviceId;
  device$;
  state;
  stateSubscription;
  desiredStateUpdated = false;
  reportedStateUpdated = false;
  editDesiredStateStatus = false;
  desiredState;

  JSON = JSON;

  constructor(private route: ActivatedRoute,
              private deviceService: DeviceService) {
  }

  ngOnInit() {
    this.route.paramMap.subscribe(params => {
      this.deviceId = params.get('deviceId');
      this.device$ = this.deviceService.getOne(this.deviceId);
      this.deviceService.getState(this.deviceId).subscribe(state => {
        if (state) {
          this.state = state;
        } else {
          this.state = {
            reported: {},
            desired: {}
          };
        }
        this.stateSubscription = this.deviceService.streamState(this.deviceId).subscribe((data) => {
          this.handleStateUpdate(data);
        });
      });
    });
  }

  private handleStateUpdate(data) {
    if (data.reportedState && data.reportedState !== null) {
      this.state.reported = data.reportedState;
      this.reportedStateUpdated = true;
      setTimeout(() => {
        this.reportedStateUpdated = false;
      }, 500);
    }
    if (data.desiredState && data.desiredState !== null) {
      this.state.desired = data.desiredState;
      this.desiredStateUpdated = true;
      setTimeout(() => {
        this.desiredStateUpdated = false;
      }, 500);
    }
  }

  private editDesiredState() {
    this.editDesiredStateStatus = true;
    this.desiredState = JSON.stringify(this.state.desired.data, null, 2);
  }

  private updateDesiredState() {
    const desiredStateObject = JSON.parse(this.desiredState);
    this.deviceService.updateDesiredState(this.deviceId, desiredStateObject)
      .subscribe(() => {
        this.deviceService.getState(this.deviceId).subscribe(state => {
          this.state = state;
          this.editDesiredStateStatus = false;
        });
      });
  }

  ngOnDestroy(): void {
    if (this.stateSubscription) {
      this.stateSubscription.unsubscribe();
    }
  }

}
