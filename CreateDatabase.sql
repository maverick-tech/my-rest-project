CREATE SCHEMA TestSchema;
GO

CREATE TABLE TestSchema.Movies (
  Id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
  Name NVARCHAR(50),
  Year NVARCHAR(50)
);
GO

INSERT INTO TestSchema.Movies (Name, Year) VALUES
(N'Avengers', N'2016'),
(N'Thor', N'2014'),
(N'Ironman', N'2010');
GO

SELECT * FROM TestSchema.Movies;
GO