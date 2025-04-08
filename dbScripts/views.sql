DROP VIEW IF EXISTS `mydb`.`UsuariosDistritosView`;

CREATE VIEW IF NOT EXISTS `UsuariosDistritosView` AS 
SELECT
    d.DistritoId,
    d.Distrito,
    c.Canton,
    p.Provincia,
    U1.Usuarios_UsuarioId as `UsuarioId`
FROM
    Distritos d
INNER JOIN Usuarios_has_Distritos U1
    ON U1.Distritos_DistritoId = d.DistritoId
INNER JOIN Cantones c
    ON d.CantonId = c.CantonId
INNER JOIN Provincias p
    ON p.ProvinciaId = c.ProvinciaId;



DROP VIEW IF EXISTS `mydb`.`ProvinciasCantonesView`;

CREATE VIEW IF NOT EXISTS `ProvinciasCantonesView` AS 
SELECT
    p.ProvinciaId AS `ProvinciaId`,
    c.CantonId AS `CantonId`,
    d.DistritoId AS `DistritoId`
FROM
    Distritos d
INNER JOIN Cantones c
    ON c.CantonId = d.CantonId
INNER JOIN Provincias p
    ON p.ProvinciaId = c.ProvinciaId
