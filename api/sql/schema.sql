CREATE TYPE platform AS ENUM ('POLYGON', 'KIS', 'BUYCYCLE');
CREATE TYPE market AS ENUM ('STOCK', 'CRYPTO', 'FOREX', 'OPTION', 'FUTURE');
CREATE TYPE locale AS ENUM ('KOR', 'USA');

/* metadata table for product data */
CREATE TABLE product_meta (
  id          varchar(24) NOT NULL,  /* an identifier form : {market}.{symbol}.{locale} */
  symbol      varchar(8)  NOT NULL, /* product symbol */
  platform    playform    NOT NULL, /* ex) polygon,   */
  locale      locale      NOT NULL, /* locale of product */
  market      market      NOT NULL, /* market of product */

  name        varchar(32), /* human readable stock name, optional */
  description text,        /* product description : can be gain by external apis, or written manually by admin, optional */

  PRIMARY KEY (id),
  UNIQUE (symbol, locale, market)
);

CREATE INDEX idx_platform_market ON product_meta (platform, market);