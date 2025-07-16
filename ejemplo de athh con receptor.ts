
//iterceptors/token.ts
import {HttpInterceptorFn}from "@angular/common/http";
export const authInterceptor: HttpInterceptorFn= {req, next} => {
    const token = localStorage.getItem("token");
    if (token) {
        req = req.clone({
            setHeaders: {
                Authentication:'Bearer ${token}'
            }
        });
    }
    return next (req);
};

import { Injectable } from '@angular/core';
import {HttpInterceptor, HttpRequest,  HttpHandler,  HttpEvent}
   from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const token = localStorage.getItem('authToken');

    if (token) {
      const authReq = req.clone({
        setHeaders: {
          Authorization: `Bearer ${token}`
        }
      });
      return next.handle(authReq);
    }

    return next.handle(req);
  }
}


/service/auth/auth.service.ts

import (Injectable) from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, BehavoirSubject, of} from "rxjs";
import {catchError, tap} from "rxjs/operators";
import {Router} from "@angular/router";

//se definen interfaces pra la repueta 
interface LoginResponse{
    token : string;
}
interface User{
    id:number;
    username:string;
    email:string;
}

@Injectable ({
    providedIn: 'root'
})
export class AuthService{
    //url de api de backend
    private apiUrl = " http://localhost:3000/api";

    //BehavoirSubject manejar estado de atenticacion de usuario y permitir qeuu otros componentes se suscriban a el 
    private _isAuthenticated = new BehavoirSubject<boolean>(false);

    private _currentUser = new BehavoirSubject<User>(null);
    currentUser$ = this._currentUser.asObservable();

    constructor(private http: HttpClient, private router: Router) {
        //al inicializar el servicio se ya ahi un token en localsorage
        this.checkAuthStatus();        
    }

    private checkAuthStatus(): void{
        const toen = this.getToken();
        if(token) {
            this._isAuthenticated.next(true);

        }
        else {
            this_isAuthenticated.next(false);
        }
    }

    @param credentials
    @returns

  

        login(credentials: {username: string; password: string}){
            this.http.post('$(this.apirl)auth/login', Credentials).subscribe({
                next: (response: any) => {
                    if (response.token) {
                        localStorage.setItem(this.tokenKey, response.token);
                        localStorage.setItem(this.tokenKey, response.user_id.toString());
                    }
                },
                error: (error) => {
                    console.error("login failed", error);
                }
            });
        }

    getUser(): Observable<User> {
        return this.http.get<User>(´${this.apiUrl}/user´).pipe(
            tap(user => {
                this._currentUser.next(user);
            }),
            catchError(error => {
                console.error("Error al obtener los datos del usuario:", error);
            return of(error);
            })
        );
    }
    
    logout(): void{
        thiis.removeTken();
        this._isAuthenticated.nex(false);
        this._currentUser.next(null);
        this.router.navigate(['/auth/login']);
    }

    @param token
    private setTokn(token: string): void{
        localStorage.setItem("authToken", token);
    }

    @return

    getToken(); string | null {
        return localStorage.getItem("authToken");
    }

    private removeToken(): void {
        localStorage.removeItem("authToken");
    }
    
}
