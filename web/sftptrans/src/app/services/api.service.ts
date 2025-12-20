import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface FileEntry {
  name: string;
  path: string;
  size: number;
  isDir: boolean;
  modTime: string;
  permissions: string;
}

export interface StatusResponse {
  connected: boolean;
  connection: string;
  downloadDir: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: {
    message: string;
    code: string;
  };
}

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(private http: HttpClient) {}

  getStatus(): Observable<ApiResponse<StatusResponse>> {
    return this.http.get<ApiResponse<StatusResponse>>('/api/status');
  }
}
