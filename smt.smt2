(declare-fun b () (Array Int (_ FloatingPoint 11 53)))
(declare-fun a () (Array Int (_ FloatingPoint 11 53)))
(declare-fun complex!1 () (Array Int (_ FloatingPoint 11 53)))
(declare-fun complex!2 () (Array Int (_ FloatingPoint 11 53)))
(declare-fun complex!3 () (Array Int (_ FloatingPoint 11 53)))
(assert
        (and
            (=
              (select complex!1 0)
              (fp.add roundNearestTiesToEven (select a 0) (select b 0))
            )
            (=
              (select complex!1 1)
              (fp.add roundNearestTiesToEven (select a 1) (select b 1))
            )
        ))
(assert
        (and
            (=
              (select complex!2 0)
              (fp.sub roundNearestTiesToEven (select a 0) (select b 0))
            )
            (=
              (select complex!2 1)
              (fp.sub roundNearestTiesToEven (select a 1) (select b 1))
            )
        ))
(assert
        (let
            (
            (a!1 (=
                   (select complex!3 0)
                   (fp.sub
                          roundNearestTiesToEven
                          (fp.mul roundNearestTiesToEven (select a 0) (select b 0))
                          (fp.mul roundNearestTiesToEven (select a 1) (select b 1))
                   )
                 ))
            (a!2 (=
                   (select complex!3 1)
                   (fp.add
                          roundNearestTiesToEven
                          (fp.mul roundNearestTiesToEven (select a 0) (select b 1))
                          (fp.mul roundNearestTiesToEven (select a 1) (select b 0))
                   )
                 ))
            )
            (and a!1 a!2)
        ))
(assert
        (let
            (
            (a!1 (and (not (fp.gt (select a 1) (select b 1))) false))
            )
            (let
                (
                (a!2 (or (and (fp.gt (select a 1) (select b 1)) false) a!1))
                )
                (let
                    (
                    (a!3 (and (not (fp.gt (select a 0) (select b 0))) a!2))
                    )
                    (or (and (fp.gt (select a 0) (select b 0)) false) a!3)
                )
            )
        ))