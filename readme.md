# Babel Router

`babel-router` is a Go implementation of (duh) a [Babel](https://docs.google.com/document/d/1B8_FC-u9iGq4RVdUB0VTxRnriBtdFCxIbqk3bhIdidU) router. A Babel router is a service that accepts messages and forwards them closer to their destination, or stores them until they can be delivered. Routers make no guarantee of message persistence and only store-and-forward on a best effort basis.
