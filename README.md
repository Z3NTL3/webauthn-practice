### WebAuthn Demo
> I will add a Docker compose setup on 26th October

This is a crude WebAuthn demo intended for mere experimenting with the APIs. One random weekend I thought myself to implement passkeys via WebAuthn, I read the WebAuthn guide & documentation, went over some of the spec details after which I studied the WebAuthn browser APIs and Go backend interface. This project is the result of that journey.


It was fun thing to do. 

#### Disclaimer
No AI was used for the software implementation. It is only used for visuals in the Figma design file.

#### Stack
- Go backend built on ``net/http`` with WebAuthn pkg 
- React Router frontend with Tailwindcss

#### Setup

Populate ``.env`` in ``./frontend`` and ``./backend`` folders respectively.

- **Reminder:**
> Go backend uses in memory storage which shouldn't be deployed nor should it be efficient for deployment.
> Remember, this is a mere demo so data will not persist between server restarts.
>
> In other words, you will need to register again to login when you restart the server.
