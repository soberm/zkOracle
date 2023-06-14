// SPDX-License-Identifier: AML
//
// Copyright 2017 Christian Reitwiessner
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// 2019 OKIMS

pragma solidity ^0.8.0;

import "./Pairing.sol";

contract SlashingVerifier {
    using Pairing for *;

    uint256 constant SNARK_SCALAR_FIELD =
        21888242871839275222246405745257275088548364400416034343698204186575808495617;
    uint256 constant PRIME_Q =
        21888242871839275222246405745257275088696311157297823662689037894645226208583;

    struct VerifyingKey {
        Pairing.G1Point alfa1;
        Pairing.G2Point beta2;
        Pairing.G2Point gamma2;
        Pairing.G2Point delta2;
        Pairing.G1Point[7] IC;
    }

    struct Proof {
        Pairing.G1Point A;
        Pairing.G2Point B;
        Pairing.G1Point C;
    }

    function verifyingKey() internal pure returns (VerifyingKey memory vk) {
        vk.alfa1 = Pairing.G1Point(
            uint256(
                6439517704249376664692906235442780318342118835482031893884060971864661789492
            ),
            uint256(
                1686185010864556271782110304746390457496278769480827626946027989164672088252
            )
        );
        vk.beta2 = Pairing.G2Point(
            [
                uint256(
                    17563244183860288416754147985535718651186036979049992571681692584988782263624
                ),
                uint256(
                    21132774184448027783844092652067569185174100983899305038648191420022755126713
                )
            ],
            [
                uint256(
                    17609690107384958896416888711696692675489778393194675953846043446512084430811
                ),
                uint256(
                    18494460895051667026542831112571511535235580304667968896627416649925983825893
                )
            ]
        );
        vk.gamma2 = Pairing.G2Point(
            [
                uint256(
                    16582775191786911343690502069092454710142536127564066274089561568968427403055
                ),
                uint256(
                    3595930612462450972422573145759732327743280853601045291648747176283559022171
                )
            ],
            [
                uint256(
                    12107835929602751336245890340991357211659502528240241054803978531501415722907
                ),
                uint256(
                    8647600739880427175597063571970468880320675254625409653764788363382565612692
                )
            ]
        );
        vk.delta2 = Pairing.G2Point(
            [
                uint256(
                    9266229895863293216412662999884900121357977156706887713859865536642616006144
                ),
                uint256(
                    17410200073591772907154742349882700198050578329426970626343371474569842885821
                )
            ],
            [
                uint256(
                    13925678285266369253314267752339014862067763841273984873997717430034992559778
                ),
                uint256(
                    233977848065755969065119879055134642721502214896095501418325137392787753302
                )
            ]
        );
        vk.IC[0] = Pairing.G1Point(
            uint256(
                4849585577078511517204739370094510587390751712784207266204541299447615261136
            ),
            uint256(
                12295636093856004209364972214653178878660368327422552276872619245528591020388
            )
        );
        vk.IC[1] = Pairing.G1Point(
            uint256(
                21246257571916732440106937772744927244877809167951166330903224213559582107898
            ),
            uint256(
                3837258044041476288755936478719254792383504730115097470590532917690537875934
            )
        );
        vk.IC[2] = Pairing.G1Point(
            uint256(
                4445653209522437155852217435813439939553533467284634395770881711203449083468
            ),
            uint256(
                19827590505412905421489040137482882005618679854268221077931711815483777475650
            )
        );
        vk.IC[3] = Pairing.G1Point(
            uint256(
                21727513592167213582284723487145842722270565244477229645635694059147709837442
            ),
            uint256(
                9050074662658605983119036076584006478704354324155813826898480286134478899629
            )
        );
        vk.IC[4] = Pairing.G1Point(
            uint256(
                10601508109101445944547044398412023506566973439744908568889124433793575874480
            ),
            uint256(
                7181732480342001843643455693749795850106532390854014887692184199389657687464
            )
        );
        vk.IC[5] = Pairing.G1Point(
            uint256(
                20084353616071149588964425217214520309210203564017868542672869309534403502873
            ),
            uint256(
                12321813138489367718126289014629967916929835867988660705260194144142190657471
            )
        );
        vk.IC[6] = Pairing.G1Point(
            uint256(
                8102791960480189763299862628598496896364737076813521680625306529379827898256
            ),
            uint256(
                368234853172378255852724485517056590940126770942350036871617683493472029443
            )
        );
    }

    /*
     * @returns Whether the proof is valid given the hardcoded verifying key
     *          above and the public inputs
     */
    function verifyProof(
        uint256[2] memory a,
        uint256[2][2] memory b,
        uint256[2] memory c,
        uint256[6] memory input
    ) public view returns (bool r) {
        Proof memory proof;
        proof.A = Pairing.G1Point(a[0], a[1]);
        proof.B = Pairing.G2Point([b[0][0], b[0][1]], [b[1][0], b[1][1]]);
        proof.C = Pairing.G1Point(c[0], c[1]);

        VerifyingKey memory vk = verifyingKey();

        // Compute the linear combination vk_x
        Pairing.G1Point memory vk_x = Pairing.G1Point(0, 0);

        // Make sure that proof.A, B, and C are each less than the prime q
        require(proof.A.X < PRIME_Q, "verifier-aX-gte-prime-q");
        require(proof.A.Y < PRIME_Q, "verifier-aY-gte-prime-q");

        require(proof.B.X[0] < PRIME_Q, "verifier-bX0-gte-prime-q");
        require(proof.B.Y[0] < PRIME_Q, "verifier-bY0-gte-prime-q");

        require(proof.B.X[1] < PRIME_Q, "verifier-bX1-gte-prime-q");
        require(proof.B.Y[1] < PRIME_Q, "verifier-bY1-gte-prime-q");

        require(proof.C.X < PRIME_Q, "verifier-cX-gte-prime-q");
        require(proof.C.Y < PRIME_Q, "verifier-cY-gte-prime-q");

        // Make sure that every input is less than the snark scalar field
        for (uint256 i = 0; i < input.length; i++) {
            require(
                input[i] < SNARK_SCALAR_FIELD,
                "verifier-gte-snark-scalar-field"
            );
            vk_x = Pairing.plus(
                vk_x,
                Pairing.scalar_mul(vk.IC[i + 1], input[i])
            );
        }

        vk_x = Pairing.plus(vk_x, vk.IC[0]);

        return
            Pairing.pairing(
                Pairing.negate(proof.A),
                proof.B,
                vk.alfa1,
                vk.beta2,
                vk_x,
                vk.gamma2,
                proof.C,
                vk.delta2
            );
    }
}
