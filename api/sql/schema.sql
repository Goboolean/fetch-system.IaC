/* metadata table for product data */
CREATE TABLE product_meta (
  id          varchar(32) NOT NULL,  /* an identifier form : {market}.{symbol}.{locale} */
  symbol      varchar(32) NOT NULL, /* stock symbol */
  locale      varchar(32) NOT NULL, /* locale of stock */
  market      varchar(32) NOT NULL, /* market of stock */

  name        varchar(32), /* human readable stock name, optional */
  description text,        /* product description : can be gain by external apis, or written manually by admin, optional */
  PRIMARY KEY (id)
);
