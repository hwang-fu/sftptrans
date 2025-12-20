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

  listRemote(path: string): Observable<ApiResponse<FileEntry[]>> {
    return this.http.get<ApiResponse<FileEntry[]>>(`/api/remote/list?path=${encodeURIComponent(path)}`);
  }

  listLocal(path: string): Observable<ApiResponse<FileEntry[]>> {
    return this.http.get<ApiResponse<FileEntry[]>>(`/api/local/list?path=${encodeURIComponent(path)}`);
  }

  mkdir(path: string): Observable<ApiResponse<null>> {
    return this.http.post<ApiResponse<null>>('/api/remote/mkdir', { path });
  }

  rename(oldPath: string, newPath: string): Observable<ApiResponse<null>> {
    return this.http.post<ApiResponse<null>>('/api/remote/rename', { oldPath, newPath });
  }

  delete(path: string): Observable<ApiResponse<null>> {
    return this.http.delete<ApiResponse<null>>(`/api/remote/delete?path=${encodeURIComponent(path)}`);
  }

  download(path: string): Observable<ApiResponse<{ localPath: string }>> {
    return this.http.get<ApiResponse<{ localPath: string }>>(`/api/remote/download?path=${encodeURIComponent(path)}`);
  }

  upload(localPath: string, remotePath: string): Observable<ApiResponse<{ remotePath: string }>> {
    const formData = new FormData();
    // Note: For browser-based upload, we'd use file input. This is simplified.
    return this.http.post<ApiResponse<{ remotePath: string }>>(
      `/api/remote/upload?path=${encodeURIComponent(remotePath)}`,
      formData
    );
  }
}
