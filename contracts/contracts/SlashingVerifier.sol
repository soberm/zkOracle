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

    uint256 constant SNARK_SCALAR_FIELD = 21888242871839275222246405745257275088548364400416034343698204186575808495617;
    uint256 constant PRIME_Q = 21888242871839275222246405745257275088696311157297823662689037894645226208583;

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
        vk.alfa1 = Pairing.G1Point(uint256(235908514540747217696550234102734792822927370188989215279248707236423443543), uint256(9459162564630235346263868952732351221128183600489855341534612403415286582237));
        vk.beta2 = Pairing.G2Point([uint256(11552633216226527004050977633436252963144865158824855858364227360057475315633), uint256(4617691669499762953439380773167223546708288755935924053809020755405705298779)], [uint256(15432779185666049812365478593583984349634051456920723581315328358842417197156), uint256(13840447874397622752694426645083371232883872598595601470085466074950284900560)]);
        vk.gamma2 = Pairing.G2Point([uint256(16805737653148279999256106789643566460198831678586513337664266618006818898661), uint256(3401479304735104352419645860646474550044738981826970782092819428180339208178)], [uint256(17643147600599063347661475412116135700931834103598903224638292504191954115945), uint256(15084966783797425503053858443135096303015617332021224332187418528753557911343)]);
        vk.delta2 = Pairing.G2Point([uint256(21784167347855289447354796192150431301009544385408127087964466227323151756487), uint256(4861587112483455803872833170244240200423291478257198542646498145395219963939)], [uint256(18488682872822422292639674179624362626847237102821751109440802405951232531088), uint256(11979398633891054829126966727345411072929754615960200868788911052238593310822)]);
        vk.IC[0] = Pairing.G1Point(uint256(6444630310688883734295296198372567805931522477817503203225152206625763918939), uint256(2266677985988175591889745310537522618319561990943032668851982234814143635803));
        vk.IC[1] = Pairing.G1Point(uint256(17196993020616208790616422359153408086742037046680508414094363149232749654964), uint256(21607219265854126399673800398604500942305788829820973857022842822481906980862));
        vk.IC[2] = Pairing.G1Point(uint256(18654297625316280796442737056644904351802435373928387426933046925417356302708), uint256(4386220326936715372311663598618426607502529490355083696998609394893216453930));
        vk.IC[3] = Pairing.G1Point(uint256(17231609154765438714837749961639337228277963639411540966205309710165180528847), uint256(7537470450224599842909642042728466414642040511431390342360664612800911109066));
        vk.IC[4] = Pairing.G1Point(uint256(7743952039720550015497959533106542762263005617892039308482071740025165609744), uint256(5051160396099968226528281973550507031055961636591044907512306593511935350737));
        vk.IC[5] = Pairing.G1Point(uint256(3875219758118574312150569446321338149685345454041751122243253028489947779934), uint256(15890213683653381332540742250279177116902256282909951517077205799777883831085));
        vk.IC[6] = Pairing.G1Point(uint256(15448612940732394068160832004453600206278029674871465537934022987470846923100), uint256(18314391648042719621044496764336194203007208800733262682518423096263353344512));
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
            require(input[i] < SNARK_SCALAR_FIELD,"verifier-gte-snark-scalar-field");
            vk_x = Pairing.plus(vk_x, Pairing.scalar_mul(vk.IC[i + 1], input[i]));
        }

        vk_x = Pairing.plus(vk_x, vk.IC[0]);

        return Pairing.pairing(
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
