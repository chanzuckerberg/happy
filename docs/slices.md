# Slices

Slices are a way to work with only a subset of a stack. For example, you might have `frontend` and `backend` slices. At times you might only want to work with the `backend`.

Slice semantics are closely related to those of `docker-compose` [profiles](https://docs.docker.com/compose/profiles/).


---
**NOTE**

We reserve the `all` slice to indicate you want to interact with "all" slices in the stack. To that end, it is important you add the `all` profile to all your services in your docker-compose configuration.

---
